import { Component, signal, inject, input, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { WorkerItemModel } from '../models/worker/worker-item';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-start-worker',
  imports: [FormsModule],
  templateUrl: './start-worker.html',
  styleUrls: [
    './start-worker.css',
    '../../assets/styles/settings-controls.css'
  ]
})
export class StartWorker {
  errorMessage = signal<string | null>(null)
  isMinimized = signal(false)
  workers = input<WorkerItemModel[]>()
  private http = inject(HttpClient);
  selected = signal(new WorkerItemModel(""))

  onSelect(event: Event){
    const selectedWorker = (event.target as HTMLSelectElement).value;
    this.selected.set(new WorkerItemModel(selectedWorker))
  }

  onStart(){

  }

  resize() {
    this.isMinimized.set(!this.isMinimized())
  }
}
