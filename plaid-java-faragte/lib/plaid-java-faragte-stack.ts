import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ecs_patterns from "aws-cdk-lib/aws-ecs-patterns";
import * as ecr from "aws-cdk-lib/aws-ecr";

export class PlaidJavaFaragteStack extends cdk.Stack {
    constructor(scope: Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // Create VPC
        const vpc = new ec2.Vpc(this, "VPC", {});

        // Create ECS Cluster
        const cluster = new ecs.Cluster(this, "Cluster", {
            vpc: vpc,
        });

        // Import ECR repository
        const repository = ecr.Repository.fromRepositoryName(
            this,
            "javaPlaidRepoID",
            "java-plaid-image-repo"
        );

        // Create Fargate service
        const fargateService =
            new ecs_patterns.ApplicationLoadBalancedFargateService(
                this,
                "FargateService",
                {
                    cluster: cluster,
                    taskImageOptions: {
                        image: ecs.ContainerImage.fromEcrRepository(repository),
                        containerPort: 8080, // Change to your Spring Boot application's port if different
                        environment: {
                            // Add any necessary environment variables
                            // 'ENV_VAR_NAME': 'ENV_VAR_VALUE',
                        },
                    },
                    memoryLimitMiB: 1024,
                    cpu: 512,
                    desiredCount: 1,
                    publicLoadBalancer: true,
                }
            );

        // Add necessary permissions if required
        // fargateService.taskDefinition.taskRole.addManagedPolicy(
        //   aws_iam.ManagedPolicy.fromAwsManagedPolicyName('policy-name')
        // );
    }
}
