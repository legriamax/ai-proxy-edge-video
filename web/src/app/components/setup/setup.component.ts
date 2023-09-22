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

 