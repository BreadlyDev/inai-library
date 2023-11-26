from django.urls import path
from .views import *

urlpatterns = [
    path("list/review/<int:book_id>/", ReviewListAPIView.as_view()),
    path("create/review/", ReviewCreateAPIView.as_view()),
    path("change/review/<int:pk>/", ReviewRetrieveUpdateDeleteAPIView.as_view()),
]
