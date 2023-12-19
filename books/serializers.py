from rest_framework import serializers
from .models import Book, Category, Subcategory


class BookSerializer(serializers.ModelSerializer):
    rating = serializers.ReadOnlyField()
    orders = serializers.ReadOnlyField()
    reviews = serializers.ReadOnlyField()
    isPossibleToOrder = serializers.BooleanField(default=True)
    inventory_number = serializers.CharField(allow_null=True)
    e_book = serializers.FileField(allow_null=True)

    def validate(self, data):
        if data["inventory_number"] is None and data["e_book"] is None:
            raise serializers.ValidationError(
                "Хотя бы одно из полей (Инвентарный номер, электронная книга) должно быть заполнено"
            )
        return data

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
