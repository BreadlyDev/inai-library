from django.urls import path
from .views import *

urlpatterns = [
    path("order/all", OrderListAPIView.as_view()),
    path("order/create", OrderCreateAPIView.as_view()),
    path("order/<int:pk>", OrderRetrieveUpdateDestroyAPIView.as_view()),
]
