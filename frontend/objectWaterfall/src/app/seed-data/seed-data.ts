import { Component, signal, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';

class Data {
  workerName = ""
  jStr = ""
  count = 1000
}

@Component({
  selector: 'app-seed-data',
  imports: [FormsModule],
  templateUrl: './seed-data.html',
  styleUrl: './seed-data.css'
})
export class SeedData {
  seedingData = signal<Data>(new Data())
  errorMessage = signal<string | null>(null)
  isLoading = signal<boolean>(false)
  isMinimized = signal(true)
  private http = inject(HttpClient);

  onSubmit() {
    this.errorMessage.set(null)
    this.isLoading.set(true)
    this.sendSettings()
  }

  private sendSettings() {
    const payload = {
      ...this.seedingData(),
      jStr: this.seedingData().jStr
    };
    this.http.post('http://localhost:8888/seed', payload).subscribe({
      next: response => {
        this.isLoading.set(false)
      },
      error: err => {
        this.errorMessage.set(err.error.error)
        this.isLoading.set(false)
      }
    });
  }

  resize() {
    this.isMinimized.set(!this.isMinimized())
  }
}
