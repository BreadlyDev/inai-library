from django.db import models
from rest_framework.exceptions import ValidationError

from main.settings import IMAGE_FOLDER, ERROR_404_IMAGE, E_BOOKS_FOLDER

LANGUAGES = (
    ("Кыргызский", "Кыргызский"),
    ("Русский", "Русский"),
    ("Английский", "Английский"),
    ("Немецкий", "Немецкий"),
)


def validate_price(phone):
    if not phone.isdigit():
        raise ValidationError("Цена должна состоять из цифр")


def validate_edition_year(edition_year):
    if not edition_year.isdigit():
        raise ValidationError("Год издания должен состоять из цифр")


class Category(models.Model):
    title = models.CharField(max_length=150)

    class Meta:
        db_table = "categories"
        verbose_name = "Category"
        verbose_name_plural = "Categories"

    def __str__(self):
        return f"Category {self.title}"


class Subcategory(models.Model):
    title = models.CharField(max_length=150)
    category = models.ForeignKey(Category, on_delete=models.CASCADE)

    class Meta:
        db_table = "subcategories"
        verbose_name = "Subcategory"
        verbose_name_plural = "Subcategories"

    def __str__(self):
        return f"Subcategory {self.title}"


class Book(models.Model):
    author = models.CharField(max_length=150)
    title = models.CharField(max_length=150)
    description = models.TextField(default="Отсутствует описание", blank=True)
    image = models.ImageField(default=ERROR_404_IMAGE, upload_to=IMAGE_FOLDER)
    e_book = models.FileField(upload_to=E_BOOKS_FOLDER, null=True, blank=True)
    category = models.ForeignKey(Category, on_delete=models.CASCADE)
    subcategory = models.ForeignKey(Subcategory, on_delete=models.CASCADE)
    inventory_number = models.CharField(max_length=150, null=True, blank=True)
    language = models.CharField(choices=LANGUAGES, max_length=150)
    edition_year = models.CharField(max_length=4, validators=[validate_edition_year])
    purchase_price = models.CharField(max_length=10, validators=[validate_price])
    purchase_time = models.DateField()
    quantity = models.PositiveIntegerField()
    isPossibleToOrder = models.BooleanField(default=True)
    rating = models.FloatField(default=0)
    orders_quantity = models.PositiveIntegerField(default=0)
    reviews_quantity = models.PositiveIntegerField(default=0)
    total_rating = models.IntegerField(default=0)
    created_time = models.DateTimeField(auto_now_add=True)

    class Meta:
        ordering = ["reviews_quantity"]
        db_table = "books"
        verbose_name = "Book"
        verbose_name_plural = "Books"

    def __str__(self):
        return f"{self.title} book with id = {self.id}"
