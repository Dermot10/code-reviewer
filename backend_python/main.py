from fastapi import FastAPI


from prometheus_client import make_asgi_app
from backend_python.handlers import review_router

app = FastAPI()
app.include_router(review_router)
app.mount("/metrics", make_asgi_app())


if __name__ == "__main__": 
    import sys
    print(sys.executable)

