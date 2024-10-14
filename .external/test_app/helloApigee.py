import uvicorn
import time
import random
from fastapi import APIRouter, Request
from fastapi.responses import JSONResponse
import os

app = APIRouter()

PORT = os.getenv("PORT", "8080")


@app.route("/{full_path:path}", ['get', 'post', 'head', 'options', 'put'])
def log_and_return_request(req: Request):

    res = {"method": req.method, "headers": dict(req.headers), "path": req.url.path, "query": req.url.query}
    print(f"Request path: {req.url.path}{req.url.path and '?'}{req.url.query}")
    print(res)

    return JSONResponse(res)


if __name__ == '__main__':
    uvicorn.run(app, host='0.0.0.0', port=int(PORT))
