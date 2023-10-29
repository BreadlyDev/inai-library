from django.urls import path
from .views import UserRegisterAPIView, UserLoginAPIView, UserLogoutAPIView

urlpatterns = [
    path('register/', UserRegisterAPIView.as_view()),
    path('login/', UserLoginAPIView.as_view()),
    path('logout/', UserLogoutAPIView.as_view())
]
