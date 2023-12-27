from rest_framework import serializers
from .models import Book, Category, Subcategory


class BookSerializer(serializers.ModelSerializer):
    rating = serializers.ReadOnlyField()
    orders = serializers.ReadOnlyField()
    reviews = serializers.ReadOnlyField()
    is_possible_to_order = serializers.BooleanField(default=True)
    e_book = serializers.FileField(allow_null=True)

    class Meta:
        model = Book
        exclude = ("created_time",)


class CategorySerializer(serializers.ModelSerializer):
    class Meta:
        model = Category
        fields = "__all__"


class SubcategorySerializer(serializers.ModelSerializer):
    class Meta:
        model = Subcategory
        fields = "__all__"
