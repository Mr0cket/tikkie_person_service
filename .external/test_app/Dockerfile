FROM python:3.12-slim

ENV PYTHONUNBUFFERED True

ENV APP_HOME /app
WORKDIR /app
COPY . ./

# Install production dependencies.
RUN pip install --no-cache-dir -r requirements.txt

CMD gunicorn --bind 0.0.0.0:${PORT} --workers 1 --worker-class uvicorn.workers.UvicornWorker --timeout 0 helloApigee:app
