import { Component } from "@angular/core";
import { NewsletterService, NewsletterSubscription } from "src/app/services/newsletter.service";

@Component({
  selector: "app-unsubscribe",
  template: `
    <div class="container" style="padding-top:65px">
      <h3>Unsubscribe</h3>
      <p>Feel free to email c0nrad@c0nrad.io if there are any complaints.</p>

      <div class="mb-3 input-group">
        <input type="email" class="form-control" placeholder="Email Address" [(ngModel)]="subscription.email" />
        <button class="btn btn-primary" (click)="subscribe(subscription)">Unsubscribe</button>
      </div>
      <div *ngIf="isSubscribed" class="text-success">Unsubscribed!</div>
      <div *ngIf="errorMessage.length != 0" class="text-danger">{{ errorMessage }}</div>
    </div>
  `,
  styles: [],
})
export class UnsubscribeComponent {
  subscription: NewsletterSubscription = {} as NewsletterSubscription;
  isSubscribed = false;
  errorMessage = "";

  constructor(private newsletterService: NewsletterService) {}

  subscribe(subscription: NewsletterSubscription) {
    this.newsletterService.unsubscribe(subscription.email).subscribe(
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
