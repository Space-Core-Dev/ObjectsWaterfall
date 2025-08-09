import { Component, input, output } from '@angular/core';
import { WorkerItemModel } from '../models/worker/worker-item';

@Component({
  selector: 'app-workers-item',
  imports: [],
  templateUrl: './workers-item.html',
  styleUrls: [
    './workers-item.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class WorkersItem {
  worker = input.required<WorkerItemModel>()
  selected = output<boolean>()

  onSelectedHandler(selected: boolean) {
    this.selected.emit(selected)
  }
}
