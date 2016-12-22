# go-project-harness
Extremely opinionated project harness based on my preferred project standards...
  * dep management
    * [glide](https://glide.sh/)
    * [npm (optional - for UI)](https://www.npmjs.com/)
  * golang dev tools
    * go vet
    * golint
  * golang loop (triggered on go file changes)
    * build a named go file
    * run the executable target with specified args
  * npm loop (run as a child process)
    * npm run watch
