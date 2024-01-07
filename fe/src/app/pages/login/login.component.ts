import { Component } from "@angular/core";

@Component({
  selector: "app-login",
  template: `
    <div class="container" style="padding-top: 65px">
      <h3>Login</h3>
      <a class="btn btn-primary" href="/login/google">Login with Google</a>
      <a class="btn btn-primary" href="/login/github">Login with Github</a>
      <!-- <h3>Hi</h3> -->
      <!-- <button class="btn btn-primary">Login with Google</button> -->
    </div>
  `,
  styles: [],
})
export class LoginComponent {}
