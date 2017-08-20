<h1 align="center">
  <a href="https://github.com/imgasm/server">imgasm</a>
</h1>

The original webapp was written in Python (Flask), focusing on providing a completely anonymous and simplistic image and video sharing service, supporting GIF to HTML5 Video formats, as well as neural style transfer.

This new version will take the project in a completely different direction. The project will be built around features that can be profitable, rather than simply storing and serving high quality image and video files on pages with little to no advertisement.

<br>

### Design decisions

Major design decisions include:

* **Complete separation of backend and frontend.** Both web and native clients will have to do *all* communication through the HTTP REST service.  
* **Go.** The server should be extremely performant. All code should be beautifully written with top-notch error handling; there should not be a need to decipher anything.

<br>

### Contributing

If you see anything that can be improved then please let me know. If you would like to be a part of the project, working on either the server or one of the clients, then feel free to hit me up at imgasm@maike.me and we can have a chat.

PS: I am *not* working on this full-time, and I prioritize quality over (development) speed, so expect slow progress :)