import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { AppController } from './controllers/app.controller';
import { SourceController } from './controllers/source';
import Source from './models/source';
import { SourceTask } from './models/sourceTask';
import CustomTypeOrmModule from './providers/typeorm.module';
import { AppService } from './services/app';
import { SourceService } from './services/source';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    CustomTypeOrmModule.forRootAsync(null, {
      entities: [Source, SourceTask],
      synchronize: true,
    }),
  ],
  controllers: [AppController, SourceController],
  providers: [AppService, SourceService],
})
export class AppModule {}
