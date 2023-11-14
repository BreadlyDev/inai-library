from rest_framework import serializers
from .models import Order


class OrderSerializer(serializers.ModelSerializer):
    class Meta:
        model = Order
        fields = ("id", "owner", "books", "status", "comment", "due_time", "created_time")


class LibrarianOrderSerializer(serializers.ModelSerializer):

    class Meta:
        model = Order
        fields = ("status", "due_time")
