import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { ConfigService } from "./config.service";

export interface Vote {
  id: string;
  ts: number;

  writeupID: string;
  userID: string;
}

@Injectable({
  providedIn: "root",
})
export class VoteService {
  constructor(private http: HttpClient, private config: ConfigService) {}

  getVotes(): Observable<Vote[]> {
    return this.http.get<Vote[]>(`${this.config.getOrigin()}/api/v1/users/me/votes`);
  }

  newVote(vote: Vote): Observable<Vote> {
    return this.http.post<Vote>(`${this.config.getOrigin()}/api/v1/votes`, vote);
  }

  deleteVote(vote: Vote): Observable<void> {
    return this.http.delete<void>(`${this.config.getOrigin()}/api/v1/votes/${vote.id}`);
  }
}
