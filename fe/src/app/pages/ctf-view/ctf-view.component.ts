import { Location } from "@angular/common";
import { Component } from "@angular/core";
import { ActivatedRoute } from "@angular/router";
import { Challenge, ChallengeService } from "src/app/services/challenge.service";
import { CTF, CTFService } from "src/app/services/ctf.service";
import { User, UserService } from "src/app/services/user.service";
import { Vote, VoteService } from "src/app/services/vote.service";
import { Writeup, WriteupService } from "src/app/services/writeup.service";

@Component({
  selector: "app-ctf-view",
  template: `<div class="container" style="padding-top: 65px">
      <button class="float-end btn-secondary btn" [routerLink]="['/ctfs', ctf.id, 'edit']" *ngIf="me.id">Edit</button>
      <h2>{{ ctf.name }}</h2>

      <!-- <div *ngFor="let category of ctf.categories" class="badge bg-secondary">{{ category }} <button tyle="button" class="btn-close" (click)="removeCategory(category)"></button></div> -->

      <div *ngFor="let category of ctf.categories">
        <h3>{{ category }}</h3>
        <table class="table table-bordered">
          <thead>
            <tr>
              <th scope="row" class="col-2">Challenge Name</th>
              <th scope="row" class="col-1">Solves</th>
              <th scope="row" class="col-3">Short Solution Description / Tags</th>
              <th scope="row" class="col-6">Writeups</th>
            </tr>
          </thead>

          <tbody>
            <tr *ngFor="let challenge of filterChallenges(category)">
              <td>
                <p *ngIf="!challenge.isEditMode">{{ challenge.name }}</p>
                <input *ngIf="challenge.isEditMode" type="text" class="form-control" placeholder="Challenge Name" [(ngModel)]="challenge.name" />
                <button *ngIf="me.isAdmin" class="btn btn-danger btn-sm" (click)="deleteChallenge(challenge)">Delete</button>
                <button *ngIf="!challenge.isEditMode && !!me.id" class="btn btn-warning btn-sm" (click)="challenge.isEditMode = true">Edit</button>
                <button *ngIf="challenge.isEditMode" class="btn btn-warning btn-sm" (click)="updateChallenge(challenge)">Save</button>
              </td>
              <td>
                <p *ngIf="!challenge.isEditMode">{{ challenge.solves }}</p>
                <input *ngIf="challenge.isEditMode" type="number" class="form-control" placeholder="Solves" [(ngModel)]="challenge.solves" />
              </td>
              <td>
                <p *ngIf="!challenge.isEditMode">{{ challenge.shortDescription }}</p>
                <input *ngIf="challenge.isEditMode" type="text" class="form-control" placeholder="Short Description" [(ngModel)]="challenge.shortDescription" />

                <!-- tags -->
                <div>
                  <span *ngFor="let tag of challenge.tags" class="badge bg-secondary me-1"
                    >{{ tag }} <button *ngIf="challenge.isEditMode" class="btn-close" (click)="removeTag(challenge, tag)"></button
                  ></span>
                </div>
                <div *ngIf="challenge.isEditMode" class="input-group mb-3">
                  <input type="text" class="form-control" placeholder="New Tag" [(ngModel)]="newTag" />
                  <button class="btn btn-secondary" (click)="addTag(challenge, newTag)">Add Tag</button>
                </div>
              </td>
              <td>
                <ul>
                  <li *ngFor="let writeup of filterWriteups(challenge.id)">
                    <i *ngIf="!hasVoted(writeup)" (click)="addVote(writeup)" class="bi bi-heart" style="  vertical-align: 0em;"></i>
                    <i *ngIf="hasVoted(writeup)" (click)="removeVote(getVote(writeup)!, writeup)" class="bi bi-heart-fill text-danger" style=""></i>

                    {{ writeup.voteCount }}

                    <a href="{{ writeup.url }}" rel="noopener noreferrer" target="_blank">{{ writeup.url }}</a>

                    <button *ngIf="me.isAdmin || writeup.submitterID == me.id" class="btn btn-danger btn-sm" (click)="deleteWriteup(writeup)">Delete</button>
                  </li>
                </ul>
                <div class="input-group" *ngIf="!!me.id">
                  <input
                    type="text"
                    class="form-control"
                    (keyup.enter)="saveWriteup(challenge, newChallengeWriteupURL[challenge.id])"
                    placeholder="Writeup URL"
                    [(ngModel)]="newChallengeWriteupURL[challenge.id]"
                  />
                  <button class="btn btn-secondary btn-sm" (click)="saveWriteup(challenge, newChallengeWriteupURL[challenge.id])">Add Writeup</button>
                </div>
              </td>
            </tr>

            <tr *ngIf="me.id">
              <td>
                <input type="text" class="form-control" placeholder="Challenge Name" (keyup.enter)="saveChallenge(category)" [(ngModel)]="newCategoryChallenge[category].name" />
              </td>
              <td>
                <input type="number" class="form-control" placeholder="Solves" [(ngModel)]="newCategoryChallenge[category].solves" />
              </td>

              <td>
                <input type="text" class="form-control" placeholder="Short Description" [(ngModel)]="newCategoryChallenge[category].shortDescription" />
              </td>

              <td>
                <button class="btn btn-primary btn-sm" (click)="saveChallenge(category)">Save New Challenge</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="input-group mb-3" *ngIf="me.id">
        <input type="text" class="form-control" placeholder="New Category" [(ngModel)]="newCategory" />
        <button class="btn btn-secondary" (click)="addCategory(newCategory)">Add Category</button>
      </div>
    </div>
    <!-- ji --> `,
  styles: [],
})
export class CtfViewComponent {
  ctfID: string = "";
  ctf: CTF = {} as CTF;
  me: User = {} as User;
  newCategory = "";
  newCategoryChallenge: { [key: string]: Challenge } = {};
  newChallengeWriteupURL: { [key: string]: string } = {};
  newTag = "";

  challenges: Challenge[] = [];
  writeups: Writeup[] = [];
  votes: Vote[] = [];

  constructor(
    private route: ActivatedRoute,
    private ctfService: CTFService,
    private location: Location,
    private userService: UserService,
    private writeupService: WriteupService,
    private voteService: VoteService,
    private challengeService: ChallengeService
  ) {}

  ngOnInit() {
    this.userService.getMe().subscribe((me) => {
      this.me = me;
    });

    this.voteService.getVotes().subscribe((votes) => {
      this.votes = votes || [];
    });

    this.route.params.subscribe((params) => {
      if (params["id"]) {
        this.ctfID = params["id"];

        this.ctfService.getCTF(params["id"]).subscribe((ctf) => {
          this.ctf = ctf;

          for (let category of this.ctf.categories) {
            this.newCategoryChallenge[category] = { ctfID: this.ctf.id, category: category } as Challenge;
          }
        });

        this.challengeService.getChallenges(params["id"]).subscribe((challenges) => {
          if (!challenges) {
            this.challenges = [];
            return;
          }
          this.challenges = challenges;
          this.newChallengeWriteupURL = {};
          for (let challenge of this.challenges) {
            this.newChallengeWriteupURL[challenge.id] = "";
          }
        });

        this.writeupService.getWriteupsForCTF(params["id"]).subscribe((writeups) => {
          this.writeups = writeups || [];
        });
      } else {
        this.ctf.submitterID = this.me.id;
      }
    });
  }

  newCTF() {
    this.location.go("/ctfs/new");
  }

  save() {
    this.ctfService.updateCTF(this.ctf).subscribe(() => {
      this.location.go(`/ctfs/${this.ctf.id}/edit`);
    });
  }

  removeCategory(category: string) {
    if (!this.ctf.categories) {
      this.ctf.categories = [];
    }
    this.ctf.categories = this.ctf.categories.filter((t) => t !== category);
  }

  addCategory(category: string) {
    if (!this.ctf.categories) {
      this.ctf.categories = [];
    }
    this.ctf.categories.push(category);
    this.newCategory = "";
    this.newCategoryChallenge[category] = { ctfID: this.ctf.id, category: category } as Challenge;
    this.save();
  }

  saveChallenge(category: string) {
    this.challengeService.newChallenge(this.newCategoryChallenge[category]).subscribe((challenge) => {
      this.challenges.push(challenge);
      this.newCategoryChallenge[category] = { ctfID: this.ctf.id, category: category } as Challenge;
    });
  }

  updateChallenge(challenge: Challenge) {
    this.challengeService.updateChallenge(challenge).subscribe(() => {
      challenge.isEditMode = false;
    });
  }

  deleteChallenge(challenge: Challenge) {
    this.challengeService.deleteChallenge(challenge).subscribe(() => {
      this.challenges = this.challenges.filter((c) => c.id !== challenge.id);
    });
  }

  deleteWriteup(writeup: Writeup) {
    this.writeupService.deleteWriteup(writeup).subscribe(() => {
      this.writeups = this.writeups.filter((w) => w.id !== writeup.id);
    });
  }

  filterChallenges(category: string): Challenge[] {
    if (this.challenges == null) {
      return [];
    }

    return this.challenges.filter((challenge) => challenge.category.toLowerCase() === category.toLowerCase()).sort((a, b) => b.solves - a.solves);
  }

  filterWriteups(challengeID: string): Writeup[] {
    if (this.writeups == null) {
      return [];
    }

    return this.writeups.filter((writeup) => writeup.challengeID === challengeID).sort((a, b) => b.voteCount - a.voteCount);
  }

  saveWriteup(challenge: Challenge, url: string) {
    this.writeupService
      .newWriteup({
        challengeID: challenge.id,
        challengeName: challenge.name,
        challengeCategory: challenge.category,
        ctfName: this.ctf.name,
        url: url,
        ctfID: this.ctf.id,
        submitterID: this.me.id,
      } as Writeup)
      .subscribe((writeup) => {
        this.writeups.push(writeup);
        this.newChallengeWriteupURL[challenge.id] = "";
      });
  }

  hasVoted(writeup: Writeup): boolean {
    if (!this.votes) {
    }
    return !!this.votes.find((v) => v.writeupID === writeup.id);
  }

  getVote(writeup: Writeup): Vote | undefined {
    return this.votes.find((v) => v.writeupID === writeup.id);
  }

  addVote(writeup: Writeup): void {
    const vote: Vote = {
      writeupID: writeup.id,
      userID: this.me.id,
    } as Vote;

    this.voteService.newVote(vote).subscribe((v) => {
      this.votes.push(v);
      writeup.voteCount++;
    });
  }

  removeVote(vote: Vote, writeup: Writeup): void {
    console.log("remove vote", vote);
    this.voteService.deleteVote(vote).subscribe((_) => {
      this.votes = this.votes.filter((v) => vote.id !== v.id);
      writeup.voteCount--;
    });
  }

  removeTag(challenge: Challenge, tag: string) {
    challenge.tags = challenge.tags.filter((t: string) => t !== tag);
    this.challengeService.updateChallenge(challenge).subscribe((c) => {
      challenge = c;
    });
  }

  addTag(challenge: Challenge, tag: string) {
    if (!challenge.tags) {
      challenge.tags = [];
    }
    challenge.tags.push(tag);
    this.newTag = "";
    this.challengeService.updateChallenge(challenge).subscribe((c) => {
      challenge = c;
    });
  }
}
