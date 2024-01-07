import { Location } from "@angular/common";
import { Component } from "@angular/core";
import { User } from "src/app/services/user.service";
import { Writeup, WriteupService } from "src/app/services/writeup.service";

@Component({
  selector: "app-newest",
  template: ` <p>new works!</p> `,
  styles: [],
})
export class NewestComponent {
  me: User = {} as User;
  writeups: Writeup[] = [];

  constructor(private writeupService: WriteupService, private location: Location) {}

  ngOnInit(): void {
    this.writeupService.getWriteups().subscribe((writeups) => {
      this.writeups = writeups;
    });
  }
}
