from django.urls import path
from .views import *

urlpatterns = [
    path("review/all", AllReviewListAPIView.as_view()),
    path("review/book/<int:book_id>", ReviewListAPIView.as_view()),
    path("review/create", ReviewCreateAPIView.as_view()),
    path("review/<int:pk>", ReviewRetrieveUpdateDeleteAPIView.as_view()),
]
