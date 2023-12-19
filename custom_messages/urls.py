from django.urls import path
from .views import *

urlpatterns = [
    path("message/all", MessageListAPIView.as_view()),
    path("message/create", MessageCreateAPIView.as_view()),
    path("message/<int:pk>", MessageRetrieveUpdateDeleteAPIView.as_view()),
]
