import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { SourceController } from './core/source.controller';
import { SourceService } from './core/source.service';

@Module({
  imports: [],
  controllers: [AppController, SourceController],
  providers: [AppService, SourceService],
})
export class AppModule {}
