from django.db import models
from core.settings import IMAGE_FOLDER, ERROR_404_IMAGE_FOLDER


class Category(models.Model):
    title = models.CharField(max_length=150)

    class Meta:
        db_table = "categories"

    def __str__(self):
        return self.title


class Book(models.Model):
    author = models.CharField(max_length=150)
    title = models.CharField(max_length=150)
    description = models.TextField(default="", blank=True)
    image = models.ImageField(default=ERROR_404_IMAGE_FOLDER, upload_to=IMAGE_FOLDER)
    category = models.ForeignKey(Category, on_delete=models.CASCADE)
    quantity = models.IntegerField()
    isPossibleToOrder = models.BooleanField(default=True)
    rating = models.FloatField(default=0)
    orders = models.IntegerField(default=0)
    reviews = models.IntegerField(default=0)
    created_time = models.DateTimeField(auto_now_add=True)

    class Meta:
        ordering = ["reviews"]
        db_table = "books"

    def __str__(self):
        return self.title
