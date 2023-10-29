from django.urls import path
from .views import *

urlpatterns = [
    path("list/message/", MessageListAPIView.as_view()),
    path("create/message/", MessageCreateAPIView.as_view()),
    path("retrieve/update/delete/message/", MessageRetrieveUpdateDeleteAPIView.as_view()),
]
