#!/usr/bin/python3
import requests

url = "http://101.43.57.52:33380/api/admin/addChallenge"

categories = ["web", "pwn", "rev", "crypto", "misc"]

for categorie in categories:
    for i in range(1,7):
        print(categorie)
        data = {
            "category": categorie,
            "title": categorie + str(i),
            "content":"for test ~~~~~~",
            "flag":"flag{youneverkonw~~}"
        }
        cookie = {
            "jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3QiLCJpc19hZG1pbiI6dHJ1ZSwidXNlcl9pZCI6MSwiaXNzIjoibXktcHJvamVjdCIsImV4cCI6MTY3ODU1MDAzNn0.8-ahxW2SH5z0TAHnmylmX5msPzq-JAY9-APTbPYdfQs"
        }
        r = requests.post(url=url,data=data,cookies=cookie)
        print(r.text)
