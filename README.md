# kuarahy-streaming
A local notification system for streamers

# Stack
React FrontEnd, Go Lang Back End, all around awesomeness 😎🚀 </br>
Vite to facilitate Go Lang and React communication

## The Problem
If you ever streamed from any service, the notification system is usually unreliable. That is related to server latency, service latency, or others. By cutting out the parts of the system that aren't interesting for a streamer, this alleviates that and simplifies the process. You won't need to login to a service for it, no bloated software, no advertisement, no platform. You use a neat little frontend browser locally installed for you to setup your notifications, and that's it.

## ToDo

[X] React Installation </br>
[ ] API Secret Testing </br>
[ ] Question: do we even need a client secret since this is all local? </br>
[X] React FrontEnd </br>
[ ] Sound Config component </br>
[X] GoReleaser Config for packaging </br>
[ ] Question: Am I using NSIS Installer for Windows? Are we releasing this for macOS in 1.0? </br>

## For Developers only
https://dev.twitch.tv/docs/api/get-started/
https://dev.twitch.tv for API and App info
https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/

npm run dev to run Vite (h + enter to open local server; localhost:5173 is mine, you can change it!)

## Developer Notes
Do NOT run npm init in the project root (unless you need a root-level package.json for scripts).
Keep node_modules confined to /frontend (it’s heavy and unnecessary elsewhere).
Go and React communicate via API (e.g., Go serves React’s built files in production).
Separation of concerns: Go (backend) and React (frontend) are decoupled.
Cleaner builds: No node_modules pollution in the Go project.
Easier deployment: React builds to static files, which Go serves.

## Build React
cd frontend && npm run build
