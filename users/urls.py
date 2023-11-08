from django.urls import path
from .views import (UserRegisterAPIView, UserLoginAPIView, UserLogoutAPIView, UserListAPIView, UserGetAPIView,
                    GroupCreateAPIView, GroupChangeAPIView)

urlpatterns = [
    path('register/', UserRegisterAPIView.as_view()),
    path('login/', UserLoginAPIView.as_view()),
    path('logout/', UserLogoutAPIView.as_view()),
    path('get/user/', UserGetAPIView.as_view()),
    path('list/user/', UserListAPIView.as_view()),
    path('create/group/', GroupCreateAPIView.as_view()),
    path('change/group/', GroupChangeAPIView.as_view()),
]
