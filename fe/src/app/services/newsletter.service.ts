import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { ConfigService } from "./config.service";
import { Observable } from "rxjs";

export interface NewsletterSubscription {
  id: string;
  ts: number;

  email: string;
}

@Injectable({
  providedIn: "root",
})
export class NewsletterService {
  constructor(private http: HttpClient, private config: ConfigService) {}

  subscribe(email: string): Observable<void> {
    return this.http.post<void>(`${this.config.getOrigin()}/api/v1/newsletter`, { email });
  }

  unsubscribe(email: string): Observable<void> {
    return this.http.post<void>(`${this.config.getOrigin()}/api/v1/newsletter/unsubscribe`, { email });
  }
}
