import { MediaMatcher } from '@angular/cdk/layout';
import { ThrowStmt } from '@angular/compiler';
import { toBase64String } from '@angular/compiler/src/output/source_map';
import { ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { NotificationsService } from 'angular2-notifications';
import { Observable, Subject } from 'rxjs';
import { ImageUpgrade } from 'src/app/models/ImageUpgrade';
import { GlobalVars } from 'src/app/models/RTSP';
import { EdgeService } from 'src/app/services/edge.service';
import { ConfirmDialogComponent } from '../shared/confirm-dialog/confirm-dialog.component';
import { WaitDialogComponent } from '../shared/wait-dialog/wait-dialog.component';

const dockerTags = ["chryscloud/chrysedgeproxy"]

@Component({
  selector: 'app-setup',
  templateUrl: './setup.component.html',
  styleUrls: ['./setup.component.scss']
})
export class SetupComponent implements OnInit, OnDestroy {

  loading:boolean = false;
  loadingMessage:string = "Please wait ... checking settings";
  title:string = ""

  mobileQuery: MediaQueryList;

  imageUpgrade = new Subject<ImageUpgrade>();
  imageUpgrade$ = this.imageUpgrade.asObservable();

  imageUpgrades:ImageUpgrade[]= [];

  expectedResponses:number;
  gotResponses:number = 0;

  private _mobileQueryListener: () => void;

  constructor(changeDetectorRef: ChangeDetectorRef, 
      media: MediaMatcher,
      private router:Router, 
      private edgeService:EdgeService,  
      private notifService:NotificationsService,
      public dialog:MatDialog) { 

    this.mobileQuery = media.matchMedia('(max-width: 600px)');
    this._mobileQueryListener = () => changeDetectorRef.detectChanges();
    this.mobileQuery.addListener(this._mobileQueryListener);

  }

  ngOnInit(): void {

    const dialogRef = this.dialog.open(WaitDialogComponent, {
      maxWidth: "400px",
      disableClose: true,
      data: {
          title: "Please wait...",
          message: "Checking for updates"}
      });

    this.initialSetup();

    this.imageUpgrade$.subscribe(data => {
      this.gotResponses += 1;
      this.imageUpgrades.push(data);
      if (this.gotResponses == this.expectedResponses) {
        // close dialog, choose which option to proceed with
        dialogRef.close();

        if (!data.has_upgrade && data.has_image) {

          // no upgrades, latest version available
          this.router.navigate(['/local/processes']);
        } else if (!data.has_upgrade && !data.has_image) {

          // no images found...initial setup
          this.title = "Initial setup. Please choose a camera type to install";
        } else if (data.has_image && data.has_upgrade) {

          // upgrade available
          this.title = "Upgrade available";
          const dialogRef = this.dialog.open(ConfirmDialogComponent, {
            maxWidt