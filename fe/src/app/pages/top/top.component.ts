import { Location } from "@angular/common";
import { Component } from "@angular/core";
import { CommentService, Comment } from "src/app/services/comment.service";
import { Seen, SeenService } from "src/app/services/seen.service";
import { User, UserService } from "src/app/services/user.service";
import { Vote, VoteService } from "src/app/services/vote.service";
import { Writeup, WriteupService } from "src/app/services/writeup.service";

@Component({
  selector: "app-top",
  template: ` <div style="padding-top: 65px" class="container">
    <div class="input-group mb-3">
      <select class="form-select" [(ngModel)]="sortBy">
        <option value="votes">Votes</option>
        <option value="newest">Newest</option>
        <option value="comments">Comments</option>
      </select>

      <select class="form-select" [(ngModel)]="categoryFilter">
        <option value="">All Categories</option>
        <option *ngFor="let category of categories" [value]="category">{{ category }}</option>
      </select>

      <select class="form-select float-end" [(ngModel)]="durationFilter">
        <option value="day">Day</option>
        <option value="week">Week</option>
        <option value="month">Month</option>
        <option value="all">All Time</option>
      </select>
    </div>

    <div class="clearfix"></div>
    <div *ngFor="let writeup of filteredWriteups()" class="card">
      <div class="card-body py-1">
        <div class="row">
          <div class="col-1">
            <i *ngIf="!hasVoted(writeup)" (click)="addVote(writeup)" class="bi bi-heart" style="font-size:1.2rem;  vertical-align: 0em;"></i>
            <i *ngIf="hasVoted(writeup)" (click)="removeVote(getVote(writeup)!, writeup)" class="bi bi-heart-fill text-danger" style="font-size: 1.2rem"></i>

            {{ writeup.voteCount }}
          </div>

          <div class="col-11">
            <div class="row pb-0">
              <div class="col-2">
                <a [routerLink]="['/ctfs', writeup.ctfID, 'edit']"> {{ writeup.ctfName }}</a>
              </div>

              <div class="col-9" [class.text-muted]="hasSeen(writeup)">
                {{ writeup.challengeCategory }} / {{ writeup.challengeName }}
                <span style="padding-left: 5px"
                  ><a [class.text-muted]="hasSeen(writeup)" [href]="writeup.url" rel="noopener noreferrer" target="_blank" (click)="addSeen(writeup)">{{ writeup.url }}</a></span
                >
              </div>

              <div class="col-1" (click)="toggleWriteup(writeup)"><i class="bi bi-chat me-1"></i>{{ writeup.commentCount || 0 }}</div>
            </div>
            <div class="row pb-0">
              <div class="col-2">
                <!-- <p>{{ writeup.body }}</p> -->
                <p>{{ writeup.ts | dateAgo }}</p>
              </div>
              <div class="col-10">
                <span *ngFor="let tag of writeup.tags" class="badge bg-secondary me-1">{{ tag }}</span>
                <i *ngIf="me.isAdmin || writeup.submitterID == me.id" class="bi bi-pencil text-muted" [routerLink]="['/writeups', writeup.id, 'edit']"></i>
              </div>
            </div>

            <div *ngIf="writeup.isExpanded">
              <div *ngFor="let comment of writeup.comments">
                <hr />
                <p class="text-muted">{{ comment.username || "Anonymous" }} - {{ comment.ts | dateAgo }}</p>
                <p>{{ comment.body }}</p>
              </div>

              <div *ngIf="me.id">
                <textarea [(ngModel)]="writeup.newComment.body" class="form-control"></textarea>
                <button (click)="addComment(writeup)" class="btn btn-primary">Add Comment</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>`,
  styles: [""],
})
export class TopComponent {
  me: User = {} as User;
  writeups: Writeup[] = [];
  votes: Vote[] = [];
  seens: Seen[] = [];

  sortBy = "votes";
  durationFilter = "week";
  categoryFilter = "";

  categories: string[] = [];

  constructor(
    private writeupService: WriteupService,
    private location: Location,
    private voteService: VoteService,
    private commentService: CommentService,
    private userService: UserService,
    private seenService: SeenService
  ) {}

  ngOnInit(): void {
    this.writeupService.getWriteups().subscribe((writeups) => {
      this.writeups = writeups || [];

      this.writeups.forEach((writeup) => {
        writeup.newComment = {} as Comment;
        writeup.newComment.body = "";
        writeup.isExpanded = false;
        writeup.comments = writeup.comments || [];

        this.categories.push(writeup.challengeCategory.toLowerCase());
      });

      this.categories = this.categories.filter((v, i, a) => a.indexOf(v) === i);
    });

    this.voteService.getVotes().subscribe((votes) => {
      this.votes = votes || [];
    });

    this.seenService.getSeens().subscribe((seens) => {
      this.seens = seens || [];
    });

    this.userService.getMe().subscribe((me) => {
      this.me = me;
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

  hasSeen(writeup: Writeup): boolean {
    return !!this.seens.find((s) => s.writeupID === writeup.id);
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

  addSeen(writeup: Writeup): void {
    const seen: Seen = {
      writeupID: writeup.id,
      userID: this.me.id,
    } as Seen;
    this.seenService.newSeen(seen).subscribe((s) => {
      this.seens.push(s);
    });
  }

  removeVote(vote: Vote, writeup: Writeup): void {
    console.log("remove vote", vote);
    this.voteService.deleteVote(vote).subscribe((_) => {
      this.votes = this.votes.filter((v) => vote.id !== v.id);
      writeup.voteCount--;
    });
  }

  getDomain(u: string): string {
    return new URL(u).hostname;
  }

  filteredWriteups(): Writeup[] {
    let writeups = this.writeups;

    if (this.durationFilter === "day") {
      writeups = writeups.filter((w) => w.ctfEndDate * 1000 > Date.now() - 1000 * 60 * 60 * 24);
    } else if (this.durationFilter === "week" || this.durationFilter === "") {
      writeups = writeups.filter((w) => w.ctfEndDate * 1000 > Date.now() - 1000 * 60 * 60 * 24 * 7);
    } else if (this.durationFilter === "month") {
      writeups = writeups.filter((w) => w.ctfEndDate * 1000 > Date.now() - 1000 * 60 * 60 * 24 * 30);
    }

    if (this.categoryFilter) {
      writeups = writeups.filter((w) => w.challengeCategory.toLowerCase() === this.categoryFilter.toLowerCase());
    }

    if (this.sortBy === "votes" || this.sortBy == "") {
      writeups = writeups.sort((a, b) => b.voteCount - a.voteCount);
    } else if (this.sortBy == "newest") {
      writeups = writeups.sort((a, b) => new Date(b.ts).getTime() - new Date(a.ts).getTime());
    } else if (this.sortBy == "comments") {
      writeups = writeups.sort((a, b) => b.commentCount - a.commentCount);
    }

    return writeups;
  }

  toggleWriteup(writeup: Writeup): void {
    writeup.isExpanded = !writeup.isExpanded;

    if (writeup.isExpanded) {
      this.commentService.getCommentsForWriteup(writeup.id).subscribe((comments) => {
        writeup.comments = comments;
      });
    }
  }

  addComment(writeup: Writeup): void {
    writeup.newComment.writeupID = writeup.id;
    writeup.newComment.userID = this.me.id;

    this.commentService.newComment(writeup.newComment).subscribe((comment) => {
      writeup.comments.push(comment);
      writeup.commentCount++;
      writeup.newComment.body = "";
    });
  }
}
