import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { ConfigService } from "./config.service";
import { Observable } from "rxjs";

export interface CTF {
  id: string;
  ts: string;

  // startDate: number;
  endDate: number;

  submitterID: string;

  name: string;
  url: string;

  challengeCount: number;
  writeupCount: number;

  categories: string[];
}

@Injectable({
  providedIn: "root",
})
export class CTFService {
  constructor(private http: HttpClient, private config: ConfigService) {}

  getCTFs(): Observable<CTF[]> {
    return this.http.get<CTF[]>(`${this.config.getOrigin()}/api/v1/ctfs`);
  }

  newCTF(ctf: CTF): Observable<CTF> {
    return this.http.post<CTF>(`${this.config.getOrigin()}/api/v1/ctfs`, ctf);
  }

  updateCTF(ctf: CTF): Observable<CTF> {
    return this.http.put<CTF>(`${this.config.getOrigin()}/api/v1/ctfs/${ctf.id}`, ctf);
  }

  getCTF(id: string): Observable<CTF> {
    return this.http.get<CTF>(`${this.config.getOrigin()}/api/v1/ctfs/${id}`);
  }
}
