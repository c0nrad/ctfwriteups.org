import { Location } from "@angular/common";
import { Component } from "@angular/core";
import { CTF, CTFService } from "src/app/services/ctf.service";
import { User, UserService } from "src/app/services/user.service";

@Component({
  selector: "app-ctf-list",
  template: `<div class="container" style="padding-top: 65px">
    <a *ngIf="me.id" class="btn btn-primary" [routerLink]="['/ctfs', 'new', 'edit']">New CTF</a>

    <table class="table">
      <thead>
        <tr>
          <th>Name</th>
          <th>Challenge Count</th>
          <th>Writeup Count</th>
          <th>End Date</th>
        </tr>
      </thead>

      <tbody>
        <tr *ngFor="let ctf of ctfs">
          <td>
            <a [routerLink]="['/ctfs', ctf.id]">{{ ctf.name }}</a>
          </td>
          <td>{{ ctf.challengeCount }}</td>
          <td>{{ ctf.writeupCount }}</td>
          <td>{{ ctf.endDate * 1000 | date }}</td>
        </tr>
      </tbody>
    </table>

    <!-- {{ ctfs }} -->
  </div>`,
  styles: [],
})
export class CtfListComponent {
  ctfs: CTF[] = [];
  me: User = {} as User;

  constructor(private ctfService: CTFService, private location: Location, private userService: UserService) {}

  ngOnInit() {
    this.ctfService.getCTFs().subscribe((ctfs) => {
      this.ctfs = ctfs;
    });

    this.userService.getMe().subscribe((me) => {
      this.me = me;
    });
  }

  newCTF() {
    this.location.go("/ctfs/new");
  }
}
