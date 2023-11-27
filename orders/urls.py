from django.urls import path
from .views import *

urlpatterns = [
    path("list/order", OrderListAPIView.as_view()),
    path("create/order", OrderCreateAPIView.as_view()),
    path("change/order/<int:pk>", OrderRetrieveUpdateDestroyAPIView.as_view()),
]
