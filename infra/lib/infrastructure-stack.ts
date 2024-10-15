import * as cdk from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as sqs from "aws-cdk-lib/aws-sqs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as docdb from "aws-cdk-lib/aws-docdb";
import * as apigw from "aws-cdk-lib/aws-apigateway";
import * as iam from "aws-cdk-lib/aws-iam";

import * as path from "path";
import { Construct } from "constructs";
import * as secretsmanager from "aws-cdk-lib/aws-secretsmanager";

export class MyCdkProjectStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // Create SQS Queue
    const queue = new sqs.Queue(this, "PersonsQueue", {
      queueName: "persons-queue",
    });

    // Create VPC for DocumentDB
    const vpc = new ec2.Vpc(this, "MyVPC", {
      maxAzs: 2,
    });

    // Create Security Group for DocumentDB
    const dbSecurityGroup = new ec2.SecurityGroup(this, "DocDBSecurityGroup", {
      vpc,
      description: "Security group for DocumentDB",
      allowAllOutbound: true,
    });

    const bastionSecurityGroup = new ec2.SecurityGroup(
      this,
      "BastionSecurityGroup",
      {
        vpc,
        description: "Security group for Bastion vm",
        allowAllOutbound: true,
      }
    );

    const bastionKey = new ec2.KeyPair(this, "bastion-key", {
      keyPairName: "bastion2",
    });

    const ec2instance = new ec2.Instance(this, "bastion", {
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.T2,
        ec2.InstanceSize.MICRO
      ),
      vpc,
      securityGroup: dbSecurityGroup,
      machineImage: ec2.MachineImage.latestAmazonLinux2(),
      vpcSubnets: {
        subnetType: ec2.SubnetType.PUBLIC,
      },
      keyPair: bastionKey,
      role: new iam.Role(this, "bastion-role", {
        assumedBy: new iam.ServicePrincipal("ec2.amazonaws.com"),
      }),
    });

    // Create Secret for DocumentDB
    const dbSecret = new secretsmanager.Secret(this, "mongoUser", {
      secretName: "mongoUser",
      description: "App Mongo User",
      generateSecretString: {
        secretStringTemplate: JSON.stringify({
          username: "mrocket",
          password: "<PASSWORD>", // Auto-generated password
        }),
        generateStringKey: "password",
        excludePunctuation: true,
        excludeCharacters: '"@/\\',
        passwordLength: 30,
      },
    });

    // DocumentDB cluster parameter group with TLS disabled (simplify access)
    const parameterGroup = new docdb.ClusterParameterGroup(
      this,
      "DocDBParamGroup",
      {
        family: "docdb5.0",
        parameters: {
          tls: "disabled",
        },
        description: "DocumentDB cluster parameter group with TLS disabled",
      }
    );

    // Create DocumentDB cluster
    const dbCluster = new docdb.DatabaseCluster(this, "DocDB", {
      masterUser: {
        username: "mrocket",
        password: dbSecret.secretValueFromJson("password"),
      },
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.T3,
        ec2.InstanceSize.MEDIUM
      ),
      vpcSubnets: {
        subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
      },
      vpc,
      securityGroup: dbSecurityGroup,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      parameterGroup: parameterGroup,
    });

    secretsmanager.Secret;

    // Create Lambda Function (service) with dependency on DocumentDB cluster
    const lambdaFunction = new lambda.DockerImageFunction(
      this,
      "PersonService",
      {
        timeout: cdk.Duration.seconds(30),
        code: lambda.DockerImageCode.fromImageAsset(
          path.join(__dirname, "../../app")
        ),
        vpc: vpc,
        securityGroups: [dbSecurityGroup],
        vpcSubnets: {
          subnetType: ec2.SubnetType.PUBLIC,
        },
        allowPublicSubnet: true,
        functionName: "personService",
        environment: {
          APP_ENVIRONMENT: "dev",
          APP_SQS_QUEUE_NAME: queue.queueName,
          APP_MONGO_SERVER: dbCluster.clusterEndpoint.socketAddress,
          APP_DB_SECRET_ARN: dbSecret.secretArn,
        },
      }
    );

    const endpoint = new apigw.LambdaRestApi(this, "PersonEndpoint", {
      handler: lambdaFunction,
      restApiName: "PersonEndpoint",
    });

    // Grant Lambda permission to send messages to SQS
    queue.grantSendMessages(lambdaFunction);

    // Grant Lambda permission to read dbSecret
    dbSecret.grantRead(lambdaFunction);

    // Allow inbound traffic on port 27017 from Lambda's security group
    dbSecurityGroup.addIngressRule(
      dbSecurityGroup,
      ec2.Port.allTraffic(),
      "Allow all traffic between resources in this security group"
    );

    bastionSecurityGroup.addIngressRule(
      ec2.Peer.ipv4(`${ec2instance.instancePublicIp}/32`),
      ec2.Port.tcp(22)
    );

    // Output the connection string
    new cdk.CfnOutput(this, "DocDBConnectionStringOutput", {
      value: dbCluster.clusterEndpoint.socketAddress,
    });

    new cdk.CfnOutput(this, "PersonEndpointOutput", {
      value: endpoint.url,
    });

    new cdk.CfnOutput(this, "BastionKeyPairName", {
      value: bastionKey.keyPairName,
    });
  }
}
