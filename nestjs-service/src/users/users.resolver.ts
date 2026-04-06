import { Resolver, Query, Args, ResolveField, Parent } from '@nestjs/graphql';
import { UsersService } from './users.service';
import { User, Event } from '../models.graphql';
import { UseGuards } from '@nestjs/common';
import { JwtAuthGuard } from '../auth/jwt-auth.guard';
import { PrismaService } from '../prisma/prisma.service';

@Resolver(() => User)
export class UsersResolver {
  constructor(
    private usersService: UsersService,
    private prisma: PrismaService,
  ) {}

  @Query(() => [User])
  @UseGuards(JwtAuthGuard)
  async getUsers() {
    return this.usersService.findAll();
  }

  @Query(() => User)
  @UseGuards(JwtAuthGuard)
  async getUser(@Args('id') id: string) {
    return this.usersService.findOne(id);
  }

  @ResolveField(() => [Event])
  async events(@Parent() user: User) {
    const records = await this.prisma.eventParticipant.findMany({
      where: { userId: user.id },
      include: { event: true },
    });
    return records.map(r => r.event);
  }
}
