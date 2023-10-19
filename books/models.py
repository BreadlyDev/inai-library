from django.db import models
from main.settings import IMAGE_FOLDER, MEDIA_ROOT


class Category(models.Model):
    title = models.CharField(max_length=150)

    def __str__(self):
        return self.title


class Book(models.Model):

    author = models.CharField(max_length=150)
    title = models.CharField(max_length=150)
    description = models.TextField(default="", blank=True)
    image = models.ImageField(default=MEDIA_ROOT+"error_404/Error404img.png", upload_to=IMAGE_FOLDER)
    category = models.ForeignKey(Category, on_delete=models.PROTECT)
    quantity = models.IntegerField()
    isPossibleToOrder = models.BooleanField(default=True)
    rating = models.FloatField(default=0)
    orders = models.IntegerField(default=0)
    reviews = models.IntegerField(default=0)
    created_time = models.DateTimeField(auto_now_add=True)

    def __str__(self):
        return self.title
