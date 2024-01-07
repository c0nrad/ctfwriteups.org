import { Injectable } from "@angular/core";
import { ConfigService } from "./config.service";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http";
import { Comment } from "./comment.service";

export interface Writeup {
  id: string;
  ts: number;

  title: string;
  url: string;
  // body: string;

  //authorID: string
  submitterID: string;

  voteCount: number;
  commentCount: number;

  ctfID: string;
  ctfName: string;
  ctfEndDate: number;

  challengeID: string;
  challengeName: string;
  challengeCategory: string;

  tags: string[];

  // UI
  isExpanded: boolean;
  comments: Comment[];
  newComment: Comment;
}

@Injectable({
  providedIn: "root",
})
export class WriteupService {
  constructor(private http: HttpClient, private config: ConfigService) {}

  getWriteups(): Observable<Writeup[]> {
    return this.http.get<Writeup[]>(`${this.config.getOrigin()}/api/v1/writeups`);
  }

  newWriteup(writeup: Writeup): Observable<Writeup> {
    return this.http.post<Writeup>(`${this.config.getOrigin()}/api/v1/writeups`, writeup);
  }

  updateWriteup(writeup: Writeup): Observable<Writeup> {
    return this.http.put<Writeup>(`${this.config.getOrigin()}/api/v1/writeups/${writeup.id}`, writeup);
  }

  getWriteup(id: string): Observable<Writeup> {
    return this.http.get<Writeup>(`${this.config.getOrigin()}/api/v1/writeups/${id}`);
  }

  getWriteupsForCTF(ctfID: string): Observable<Writeup[]> {
    return this.http.get<Writeup[]>(`${this.config.getOrigin()}/api/v1/ctfs/${ctfID}/writeups`);
  }

  deleteWriteup(writeup: Writeup): Observable<void> {
    return this.http.delete<void>(`${this.config.getOrigin()}/api/v1/writeups/${writeup.id}`);
  }
}
