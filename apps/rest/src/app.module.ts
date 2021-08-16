import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import Source from './models/source';
import { SourceTask } from './models/sourceTask';
import CustomTypeOrmModule from './providers/typeorm.module';

@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    CustomTypeOrmModule.forRootAsync(null, {
      entities: [Source, SourceTask],
      // FIXME: using db migration instead of synchronize
      synchronize: true,
    }),
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
