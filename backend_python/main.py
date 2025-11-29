from fastapi import FastAPI
from prometheus_client import make_asgi_app
from backend_python.handlers import review_router

app = FastAPI()
app.include_router(review_router)
app.mount("/metrics", make_asgi_app())


if __name__ == "__main__": 
    
    print("Hello World")



# TODO -

# Done: preprocessing - chunking 

# aggregator

# add aggregator to service code, may require chunking 

# postprocessing
    # json results
    # file postprocessing to package as correct code file

# testing 
    # dummy results
    # go api



