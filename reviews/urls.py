from django.urls import path
from .views import *

urlpatterns = [
    path("list/review/", ReviewListAPIView.as_view()),
    path("create/review/", ReviewCreateAPIView.as_view()),
    path("retrieve/update/delete/review/<int:pk>", ReviewRetrieveUpdateDeleteAPIView.as_view()),
]
