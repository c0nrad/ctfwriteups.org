import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { ConfigService } from "./config.service";

export interface Comment {
  id: string;
  ts: number;

  userID: string;
  username: string;

  writeupID: string;
  parentCommentID: string;

  body: string;
  votes: number;
}

@Injectable({
  providedIn: "root",
})
export class CommentService {
  constructor(private http: HttpClient, private config: ConfigService) {}

  getComments(): Observable<Comment[]> {
    return this.http.get<Comment[]>(`${this.config.getOrigin()}/api/v1/users/me/comments`);
  }

  newComment(comment: Comment): Observable<Comment> {
    return this.http.post<Comment>(`${this.config.getOrigin()}/api/v1/comments`, comment);
  }

  deleteComment(comment: Comment): Observable<void> {
    return this.http.delete<void>(`${this.config.getOrigin()}/api/v1/comments/${comment.id}`);
  }

  getCommentsForWriteup(writeupID: string): Observable<Comment[]> {
    return this.http.get<Comment[]>(`${this.config.getOrigin()}/api/v1/writeups/${writeupID}/comments`);
  }
}
