import { Location } from "@angular/common";
import { Component } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { Challenge, ChallengeService } from "src/app/services/challenge.service";
import { CTF, CTFService } from "src/app/services/ctf.service";
import { User, UserService } from "src/app/services/user.service";
import { Vote, VoteService } from "src/app/services/vote.service";
import { Writeup, WriteupService } from "src/app/services/writeup.service";

@Component({
  selector: "app-ctf-edit",
  template: `<div class="container" style="padding-top: 65px">
      <div class="alert alert-info alert-dismissible fade show" role="alert" *ngIf="ctfID == 'new'">
        <strong>New CTF Instructions</strong> Thanks for adding a CTF! Please make sure the CTF doesn't already exist. You can edit the CTF after creating it. The CTF enddate is important since we
        only want to show active writeups. Thanks! If there's any question join the <a href="https://discord.gg/srDFwyHmVS">discord.</a>.
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
      </div>

      <div class="mb-3" *ngIf="me">
        <div class="form-group">
          <label class="form-label">CTF Name</label>
          <input type="text" class="form-control" placeholder="CTF Name" [(ngModel)]="ctf.name" />
        </div>

        <div class="form-group">
          <label class="form-label">CTF URL</label>
          <input type="text" class="form-control" placeholder="CTF URL" [(ngModel)]="ctf.url" />
        </div>

        <div class="form-group">
          <label>CTF End (Unix Time)</label>
          <input type="number" class="form-control" placeholder="CTF End Date" [(ngModel)]="ctf.endDate" />
          <small class="form-text text-muted">{{ ctf.endDate * 1000 | date : "medium" }}</small>
        </div>

        <button *ngIf="me" (click)="save()" class="btn btn-primary">Save</button>
      </div>

      <!-- <div *ngFor="let category of ctf.categories" class="badge bg-secondary">{{ category }} <button tyle="button" class="btn-close" (click)="removeCategory(category)"></button></div> -->
    </div>
    <!-- ji --> `,
  styles: [],
})
export class CtfEditComponent {
  ctfID: string = "";
  ctf: CTF = {} as CTF;
  me: User = {} as User;

  constructor(private route: ActivatedRoute, private ctfService: CTFService, private location: Location, private userService: UserService) {}

  ngOnInit() {
    this.userService.getMe().subscribe((me) => {
      this.me = me;
    });

    this.route.params.subscribe((params) => {
      if (params["id"]) {
        this.ctfID = params["id"];
        console.log("ctfID", this.ctfID);
        if (this.ctfID == "new") {
          console.log("new ctf");
          this.ctf = {} as CTF;
          this.ctf.submitterID = this.me.id;
          this.ctf.categories = [];
          this.ctf.endDate = Math.floor(Date.now() / 1000);
          console.log(this.ctf.endDate);
          return;
        }

        this.ctfService.getCTF(params["id"]).subscribe((ctf) => {
          this.ctf = ctf;
        });
      }
    });
  }

  save() {
    if (this.ctf.id) {
      this.ctfService.updateCTF(this.ctf).subscribe(() => {
        this.location.go(`/ctfs/${this.ctf.id}`);
        window.location.reload();
      });
    } else {
      this.ctfService.newCTF(this.ctf).subscribe((ctf) => {
        this.ctf = ctf;
        this.location.go(`/ctfs/${this.ctf.id}`);
        window.location.reload();
      });
    }
  }
}
