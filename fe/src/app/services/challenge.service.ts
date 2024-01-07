import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { ConfigService } from "./config.service";

export interface Challenge {
  id: string;
  ts: number;

  ctfID: string;
  submitterID: string;

  name: string;
  category: string;
  solves: number;

  shortDescription: string;
  tags: string[];

  // virtual
  isEditMode: boolean;
}

@Injectable({
  providedIn: "root",
})
export class ChallengeService {
  constructor(private http: HttpClient, private config: ConfigService) {}

  getChallenges(ctfID: string): Observable<Challenge[]> {
    return this.http.get<Challenge[]>(`${this.config.getOrigin()}/api/v1/ctfs/${ctfID}/challenges`);
  }

  newChallenge(challenge: Challenge): Observable<Challenge> {
    return this.http.post<Challenge>(`${this.config.getOrigin()}/api/v1/challenges`, challenge);
  }

  updateChallenge(challenge: Challenge): Observable<Challenge> {
    return this.http.put<Challenge>(`${this.config.getOrigin()}/api/v1/challenges/${challenge.id}`, challenge);
  }

  getChallenge(ctfID: string, id: string): Observable<Challenge> {
    return this.http.get<Challenge>(`${this.config.getOrigin()}/api/v1/challenges/${id}`);
  }

  deleteChallenge(challenge: Challenge): Observable<void> {
    return this.http.delete<void>(`${this.config.getOrigin()}/api/v1/challenges/${challenge.id}`);
  }
}
