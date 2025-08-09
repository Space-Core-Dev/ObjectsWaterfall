import { Component, input } from '@angular/core';
import { WorkerItemModel } from '../models/worker/worker-item';
import { WorkersItem } from '../workers-item/workers-item';

@Component({
  selector: 'app-workers-list',
  imports: [WorkersItem],
  templateUrl: './workers-list.html',
  styleUrl: './workers-list.css'
})
export class WorkersList {
  workers = input<WorkerItemModel[]>()
}
