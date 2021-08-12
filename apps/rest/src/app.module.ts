import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { SourceController } from './core/source.controller';

@Module({
  imports: [],
  controllers: [AppController, SourceController],
  providers: [AppService],
})
export class AppModule {}
