import { Component } from "@angular/core";
import { User, UserService } from "src/app/services/user.service";

@Component({
  selector: "app-header",
  template: `
    <header data-bs-theme="dark">
      <nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark">
        <div class="container-fluid">
          <a class="navbar-brand" [routerLink]="['/']">
            <img src="/assets/logo.png" alt="Logo" width="30" height="24" class="d-inline-block align-text-top" />

            CTFWriteups.org</a
          >
          <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarCollapse" aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
          </button>
          <div class="collapse navbar-collapse" id="navbarCollapse">
            <ul class="navbar-nav me-auto mb-2 mb-md-0">
              <li class="nav-item">
                <a class="nav-link" [routerLinkActive]="'active'" [routerLinkActiveOptions]="{ exact: true }" [routerLink]="['/']">Home</a>
              </li>
              <!-- <li class="nav-item">
                <a class="nav-link" [routerLinkActive]="'active'" [routerLink]="['/newest']">Newest</a>
              </li> -->
              <li class="nav-item">
                <a class="nav-link" [routerLinkActive]="'active'" [routerLink]="['/ctfs']">CTFs</a>
              </li>
              <li class="nav-item">
                <a class="nav-link disabled" [routerLinkActive]="'active'" [routerLink]="['/search']">Search (coming soon)</a>
              </li>
              <li class="nav-item">
                <a class="nav-link" [routerLinkActive]="'active'" [routerLink]="['/newsletter']">Weekly Email Digest</a>
              </li>

              <li class="nav-item">
                <a class="nav-link" [routerLinkActive]="'active'" [routerLink]="['/about']">About</a>
              </li>
            </ul>

            <div class="col-md-3 text-end">
              <button *ngIf="!me.id" type="button" class="btn btn-warning" [routerLink]="['/login']">Login</button>
              <a *ngIf="me.id" type="button" class="btn btn-warning" href="/logout">Logout</a>
            </div>
          </div>
        </div>
      </nav>
    </header>
  `,
  styles: [],
})
export class HeaderComponent {
  me: User = {} as User;

  constructor(private userService: UserService) {}

  ngOnInit() {
    this.userService.getMe().subscribe((me) => {
      this.me = me;
    });
  }

  logout() {
    this.userService.logout().subscribe(() => {
      window.location.href = "/";
    });
  }
}
