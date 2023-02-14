# Goptcha
Proof of consept captcha solution done in Fiber and Golang.

localhost:3000/captcha returns an randomly generated image, which contains text which is passed to the localhost:3000/checker endpoint for verification.


![image](https://user-images.githubusercontent.com/61495413/218850589-9e30b6dd-4f69-4260-83fc-809644e5e6db.png)


## localhost:3000/checker
```JSON
{"captcha": "BQYLXCTC"}
```
