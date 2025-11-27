from fastapi import FastAPI
# from prometheus_client import make_asgi_app
from handlers import router as review_router


app = FastAPI()
app.include_router(review_router)
# app.mount("/metrics", make_asgi_app())
