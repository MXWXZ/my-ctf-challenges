import random

import face_recognition
import numpy as np
from PIL import Image

# std_image = face_recognition.load_image_file("std.jpg")
unknown_image = face_recognition.load_image_file("test.jpg")

std_encoding=[
 -8.69139656e-02,  8.30148682e-02 , 1.45035293e-02, -1.27609253e-01,
 -1.42700657e-01, -1.58593412e-02 ,-9.87722948e-02, -1.23219922e-01,
  1.22708268e-01, -1.35270610e-01 , 2.30035380e-01, -1.23880222e-01,
 -1.93354771e-01, -8.94580930e-02 ,-7.93846995e-02,  2.35654935e-01,
 -1.81906566e-01, -1.34962142e-01 ,-1.31788421e-02, -1.04968855e-02,
  4.10739481e-02,  2.44885264e-03 , 8.52121785e-03,  5.79290688e-02,
 -1.15343466e-01, -3.23355764e-01 ,-8.69766697e-02, -2.12586801e-02,
 -9.11531225e-02, -3.72300223e-02 ,-2.80866250e-02,  1.02462806e-01,
 -1.71462923e-01, -2.73887850e-02 , 4.65847105e-02,  6.94189966e-02,
  2.20984984e-02, -8.01130161e-02 , 1.72256276e-01,  1.52742490e-04,
 -2.54432797e-01,  5.17657027e-02 , 1.13474540e-01,  2.19928578e-01,
  1.68304369e-01,  1.28403883e-02 ,-1.04458071e-02, -1.59635231e-01,
  1.74563184e-01, -1.74656272e-01 , 1.19449571e-04,  1.32924736e-01,
  4.52756137e-02, -5.11706285e-02 , 1.84679162e-02, -7.74622187e-02,
  2.99685597e-02,  1.66548729e-01 ,-1.57246217e-01, -3.03353313e-02,
  9.47528481e-02, -6.63631782e-02 ,-3.17470208e-02, -1.85560584e-01,
  2.26004064e-01,  1.28806546e-01 ,-1.15559876e-01, -2.06283614e-01,
  1.40707687e-01, -1.00104943e-01 ,-8.33150819e-02,  8.25207531e-02,
 -1.33005619e-01, -1.90996230e-01 ,-2.95138747e-01, -2.70678457e-02,
  3.30062211e-01,  1.28746748e-01 ,-1.88333243e-01,  5.84503338e-02,
 -8.36766977e-03, -7.47905578e-03 , 1.23152651e-01,  1.65390745e-01,
  5.01543283e-03,  1.08317155e-02 ,-8.22547823e-02, -4.03350629e-02,
  2.58023173e-01, -4.20480780e-02 ,-2.24346798e-02,  2.48134851e-01,
 -5.13138250e-04,  6.34072348e-02 , 6.94152107e-03, -9.12788417e-03,
 -1.11195974e-01,  3.06070670e-02 ,-1.62505597e-01, -1.20745702e-02,
 -1.50425863e-02, -1.41657144e-02 ,-1.81038231e-02,  1.26067802e-01,
 -1.41881093e-01,  1.04972236e-01 ,-5.23118973e-02,  3.43461856e-02,
 -2.61395201e-02, -2.75162887e-02 ,-2.53709070e-02, -3.63143757e-02,
  1.08865552e-01, -2.02156767e-01 , 1.07431002e-01,  8.50366130e-02,
  7.95102417e-02,  1.08320944e-01 , 1.53148308e-01,  8.43793526e-02,
 -2.67507583e-02, -3.10356300e-02 ,-2.16474622e-01, -2.27650702e-02,
  1.20539531e-01, -9.48047191e-02 , 1.40443712e-01,  5.64389490e-03,
]

# cmp_encoding = face_recognition.face_encodings(unknown_image)[0]
# print(face_recognition.face_distance([std_encoding], cmp_encoding))
# exit()

results=1
step=0
top, right, bottom, left = face_recognition.face_locations(unknown_image)[0]
for i in range(top, bottom, 23):
    for j in range(left, right, 23):
        step+=1
        print(step,end='\r')
        save = np.copy(unknown_image)

        unknown_image[i][j] = [random.randrange(0,256), random.randrange(0,256), random.randrange(0,256)]
        unknown_image[i][j+1] = [random.randrange(0,256), random.randrange(0,256), random.randrange(0,256)]
        unknown_image[i+1][j] = [random.randrange(0,256), random.randrange(0,256), random.randrange(0,256)]
        unknown_image[i+1][j+1] = [random.randrange(
            0, 256), random.randrange(0, 256), random.randrange(0, 256)]
        try:
            cmp_encoding = face_recognition.face_encodings(unknown_image)[0]
            tmp = face_recognition.face_distance([std_encoding], cmp_encoding)
        except Exception as e:
            unknown_image = np.copy(save)
            continue
        if tmp<results:
            print(step,tmp)
            results=tmp
            pil_image = Image.fromarray(unknown_image)
            pil_image.save('out.jpeg', quality=100, subsampling=0)
        else:
            unknown_image = np.copy(save)
