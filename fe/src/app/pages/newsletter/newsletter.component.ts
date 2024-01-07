import { Component } from "@angular/core";
import { NewsletterService, NewsletterSubscription } from "src/app/services/newsletter.service";

@Component({
  selector: "app-newsletter",
  template: `
    <div class="container" style="padding-top:65px">
      <h3>Newsletter</h3>
      <p>Weekly email digest of the top voted CTF writeups + c0nrad's favorites.</p>
      <p>Sent out weekly on Wednesdays.</p>

      <div class="mb-3 input-group">
        <input type="email" class="form-control" placeholder="Email Address" [(ngModel)]="subscription.email" />
        <button class="btn btn-primary" (click)="subscribe(subscription)">Subscribe</button>
      </div>
      <div *ngIf="isSubscribed" class="text-success">Subscribed!</div>
      <div *ngIf="errorMessage.length != 0" class="text-danger">{{ errorMessage }}</div>
    </div>
  `,
  styles: [],
})
export class NewsletterComponent {
  subscription: NewsletterSubscription = {} as NewsletterSubscription;
  isSubscribed = false;
  errorMessage = "";

  constructor(private newsletterService: NewsletterService) {}

  subscribe(subscription: NewsletterSubscription) {
    this.newsletterService.subscribe(subscription.email).subscribe(
      (subscription) => {
        this.isSubscribed = true;
      },
      (err) => {
        // this.isSubscribed = true;
        this.errorMessage = err.error;
      }
    );
  }
}
