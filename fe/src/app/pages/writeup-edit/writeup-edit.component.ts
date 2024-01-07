import { Location } from "@angular/common";
import { Component } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { User, UserService } from "src/app/services/user.service";
import { Writeup, WriteupService } from "src/app/services/writeup.service";

@Component({
  selector: "app-writeup-edit",
  template: `
    <div class="container" style="padding-top:65px">
      <!-- <h3>{{ writeup.title }}</h3> -->

      <a [routerLink]="['/ctfs', writeup.ctfID, 'edit']">{{ writeup.ctfName }}</a> / {{ writeup.challengeCategory }} / {{ writeup.challengeName }}

      <div class="mb-3">
        <label class="form-label">URL</label>
        <input type="text" class="form-control" placeholder="URL" [(ngModel)]="writeup.url" [disabled]="!me.isAdmin" />
      </div>

      <div class="mb-3">
        <label class="form-label">Tags</label>

        <div class="input-group">
          <input type="text" class="form-control" placeholder="Tags" [(ngModel)]="newTag" (keyup.enter)="addTag(newTag)" />
          <button class="btn btn-secondary" (click)="addTag(newTag)">Add Tag</button>
        </div>

        <div *ngFor="let tag of writeup.tags" class="badge bg-secondary me-1">{{ tag }} <button tyle="button" class="btn-close" (click)="removeTag(tag)"></button></div>

        <div class="form-text">Tags are used for searching. "XSS", "nodejs", "ROP", etc.</div>
      </div>

      <button *ngIf="me.isAdmin" class="btn btn-primary" (click)="save()">Save</button>
    </div>
  `,
  styles: [],
})
export class WriteupEditComponent {
  writeup: Writeup = {} as Writeup;
  newTag = "";

  me: User = {} as User;

  constructor(private route: ActivatedRoute, private writeupService: WriteupService, private location: Location, private userService: UserService) {}

  ngOnInit() {
    this.userService.getMe().subscribe((me) => {
      this.me = me;

      this.route.params.subscribe((params) => {
        if (params["id"]) {
          this.writeupService.getWriteup(params["id"]).subscribe((writeup) => {
            this.writeup = writeup;
          });
        } else {
          this.writeup.submitterID = this.me.id;
        }
      });
    });
  }

  removeTag(tag: string) {
    if (!this.writeup.tags) {
      this.writeup.tags = [];
    }
    this.writeup.tags = this.writeup.tags.filter((t) => t !== tag);
    this.save();
  }

  addTag(tag: string) {
    if (!this.writeup.tags) {
      this.writeup.tags = [];
    }
    this.writeup.tags.push(tag);
    this.newTag = "";
    this.save();
  }

  save() {
    if (this.writeup.id) {
      this.writeupService.updateWriteup(this.writeup).subscribe(() => {
        this.location.go(`/writeups/${this.writeup.id}/edit`);
      });
    } else {
      this.writeupService.newWriteup(this.writeup).subscribe(() => {
        this.location.go(`/writeups/${this.writeup.id}/edit`);
      });
    }
  }
}
