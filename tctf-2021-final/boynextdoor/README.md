# TCTF2021-Final boynextdoor

The aim of the challenge is to break a well known [face recognition library](https://github.com/ageitgey/face_recognition) in Github, that is, given a target vector, find some face to best match the vector.

## GAN? FGSM?

You may refer to Generative Adversarial Network or Gradient Sign Method to break it in AI methods. But this may take you a lot of time to analysis the model and learn the detail of the library which is not the aim of this challenge, you can do it but need longer time.

## Hill Climbing

A good method is using hill climbing approach, randomly change some bit of the image and see whether the results get better, repeat until we get the best one. However there are some problem we need to solve:

- local optimum: The hill climbing methods could only get local optimum, not global optimum.
- search space: The search space of the image may be so large that we can't make it in reasonable time, we need some methods to reduce it.

## Local Optimum

If we choose wrong picture, we may lost in local optimum and can't get the best results no matter what we do, so in order to increase our chance and reduce training time, we need to find the best match face as our initial face. In this official writeup, I use the 1 million face in [Kaggle](https://www.kaggle.com/tunguz/1-million-fake-faces), note that I don't use all of them, just a small part.

You can use [dataset.py](dataset.py) to find the best match face. I find it in 10min and get best face at 28% distance([test.jpeg](test.jpeg)).

## Search Space

See [exp.py](exp.py), I use some methods to reduce the search space, including batch update and step update, you can also minimize it by reducing the image dimensions and so on.

However, you can not over reduce the search space so that you won't have enough space to minimize your result.

## Result

We can generate a 24% match face in 535 steps within 10 min([535.jpeg](535.jpeg))

## Notes

- This library only count face region to generate vector, so bigger face in the image could give us more chance to match.
- PIL seems to lose the image quality when save as jpeg, in the exploit we just ignore it for it's no more than 1% typically.